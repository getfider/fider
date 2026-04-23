# Jenkins CI/CD Setup for Fider on Windows

This guide provides step-by-step instructions to configure and run the Fider CI/CD pipeline on Jenkins running on Windows.

## Prerequisites

### 1. Jenkins Installation
- **Jenkins Version**: 2.400+
- **Java**: OpenJDK 11+ or Oracle JDK 11+
- **OS**: Windows Server 2016 or later / Windows 10/11 Pro

### 2. Windows Agent/Node Setup

Jenkins requires the following tools to be installed on the Windows agent/node:

#### Go
- **Version**: 1.25+
- **Download**: https://go.dev/dl/
- **Install**: Extract to `C:\Go` or add to PATH
- **Verify**: Run `go version` in PowerShell

#### Node.js
- **Version**: 22.x LTS or 22.x Current
- **Download**: https://nodejs.org/
- **Install**: Use Windows Installer
- **Verify**: Run `node --version` and `npm --version` in PowerShell

#### Docker (Optional, for Docker image builds)
- **Version**: Latest (Docker Desktop for Windows)
- **Download**: https://www.docker.com/products/docker-desktop
- **Note**: Requires Windows 10/11 Pro with WSL 2 enabled
- **Verify**: Run `docker --version` in PowerShell

#### PostgreSQL (Optional, for local testing)
- **Version**: 17+
- **Download**: https://www.postgresql.org/download/windows/
- **Note**: Can also use Docker containers for testing

### 3. Required Jenkins Plugins

Install these plugins via **Manage Jenkins** → **Manage Plugins**:

1. **Pipeline** (Declarative Pipeline)
   - ID: `workflow-aggregator`
   
2. **Git**
   - ID: `git`
   
3. **GitHub Integration**
   - ID: `github-api` and `github`
   
4. **Blue Ocean** (Optional, for better UI)
   - ID: `blueocean`
   
5. **Pipeline: GitHub Groovy Libraries**
   - ID: `pipeline-github-lib`
   
6. **Email Extension**
   - ID: `email-ext`
   
7. **HTML Publisher**
   - ID: `htmlpublisher`
   
8. **Code Coverage**
   - ID: `cobertura` or `jacoco`

## Configuration Steps

### Step 1: Create New Pipeline Job

1. In Jenkins, click **New Item**
2. Enter job name: `Fider-CI-Pipeline`
3. Select **Pipeline**
4. Click **OK**

### Step 2: Configure Pipeline

1. Go to **Pipeline** section
2. Select **Pipeline script from SCM**
3. Set **SCM** to **Git**

#### Git Configuration
- **Repository URL**: `https://github.com/your-org/fider.git`
- **Credentials**: Add your GitHub credentials
- **Branches to build**: `*/main` and `*/develop`
- **Script Path**: `Jenkinsfile`

### Step 3: Configure Build Triggers

1. Go to **Build Triggers** section
2. Enable one or more:
   - **GitHub hook trigger for GITScm polling** (requires GitHub webhook)
   - **Poll SCM** with schedule: `H * * * *` (hourly)

### Step 4: Setup GitHub Integration (Optional)

1. In Jenkins, go to **Manage Jenkins** → **Configure System**
2. Find **GitHub** section
3. Add your GitHub server credentials
4. Test connection

To add webhook in GitHub:
1. Go to Repository Settings → Webhooks
2. Click **Add webhook**
3. Payload URL: `https://your-jenkins-url/github-webhook/`
4. Content type: `application/json`
5. Events: Push events and Pull requests

### Step 5: Configure Global Tools

1. Go to **Manage Jenkins** → **Global Tool Configuration**

#### Go Configuration
- Name: `Go-1.25`
- Go version: `1.25`
- Install automatically: ✓ (optional)

#### Node.js Configuration
- Name: `Node-22`
- Node.js version: `22.x`
- Install automatically: ✓ (optional)

## Running the Pipeline

### Method 1: Manual Trigger
1. Open your pipeline job in Jenkins
2. Click **Build Now**
3. Monitor the build in real-time

### Method 2: GitHub Webhook (Recommended)
1. Push changes to `main` or configured branch
2. Jenkins automatically triggers build via webhook
3. Monitor build status in GitHub PR/commit

### Method 3: Scheduled
- Pipeline runs on schedule defined in **Poll SCM** trigger

## Pipeline Stages

The pipeline executes these stages sequentially:

### 1. Checkout
- Clones repository code
- Duration: ~10 seconds

### 2. Setup Environment
- Validates Go, Node.js, npm installation
- Creates necessary directories
- Duration: ~5 seconds

### 3. Test UI (Frontend)
- Installs npm dependencies
- Runs ESLint checks
- Runs Jest tests with coverage
- Duration: ~2-3 minutes

### 4. Test Server (Backend)
- Installs Go linter (golangci-lint)
- Runs linting checks
- Runs Go tests with coverage
- Duration: ~3-5 minutes

### 5. Build Docker Image
- Builds Docker image (if Docker available)
- Tags image with commit SHA
- Duration: ~2-3 minutes

### 6. E2E Tests (Optional)
- Installs Playwright
- Runs end-to-end tests
- Duration: ~5-10 minutes

### Total Pipeline Duration: 15-30 minutes

## Post-Build Actions

After pipeline completes, the following artifacts are archived:

- **Frontend Coverage**: `coverage/lcov-report/`
- **Backend Coverage**: `coverage-report.html`
- **Lint Reports**: `eslint-report.json`, `golangci-report.txt`
- **Build Logs**: `docker-build.log`
- **E2E Results**: `e2e-test.log`, `playwright-report/`
- **Summary Reports**: `CI_REPORT.md`, `*-test-summary.md`

### Accessing Artifacts

1. Open the Jenkins build
2. Click **Build Artifacts**
3. Download individual files or entire directories

### HTML Reports

1. After build completes, scroll to **HTML Reports** section
2. View interactive reports:
   - Code Coverage - Frontend
   - Code Coverage - Backend

## Environment Variables

The pipeline uses these environment variables:

```
GO_VERSION = '1.25'
NODE_VERSION = '22.x'
POSTGRES_VERSION = '17'
GOPATH = ${WORKSPACE}\go
PATH = ${WORKSPACE}\go\bin;${PATH}
```

You can override these in Jenkins job configuration under **Build Environment** → **Add timestamps** and custom environment variables.

## Troubleshooting

### Issue: "Command not found" for Go or Node.js

**Solution:**
1. Verify installation: Run in PowerShell
   ```powershell
   go version
   node --version
   npm --version
   ```
2. If not found, add to PATH:
   - Go: `C:\Go\bin`
   - Node.js: `C:\Program Files\nodejs`

### Issue: Docker build fails

**Solution:**
1. Verify Docker is running: `docker --version`
2. If Docker Desktop not running, start it
3. Check Docker resources (CPU/Memory)
4. For WSL 2 issues, verify WSL installation

### Issue: Tests fail with "npm: command not found"

**Solution:**
1. Restart Jenkins service to refresh PATH
2. Or explicitly use `C:\Program Files\nodejs\npm`
3. Check npm installation: `npm --version`

### Issue: "Permission denied" on Linux agents (if used)

**Solution:**
- This pipeline is Windows-specific, use Unix agents for Linux

### Issue: Pipeline timeout

**Solution:**
1. Increase timeout in Jenkinsfile: Change `timeout(time: 1, unit: 'HOURS')` to higher value
2. Check for long-running tests
3. Enable parallel execution if possible

## Advanced Configuration

### Notifications

Add email notifications by modifying the post section:

```groovy
post {
    always {
        // ... existing code ...
        emailext(
            subject: "Fider CI Build #${BUILD_NUMBER}",
            body: "Check console output at ${BUILD_URL}",
            to: "team@example.com",
            recipientProviders: [developers(), requestor()]
        )
    }
}
```

### Artifact Retention

Modify artifact archival to keep fewer builds:

```groovy
buildDiscarder(logRotator(numToKeepStr: '10'))  // Keep last 10 builds
```

### Parallel Execution

Split frontend and backend tests to run in parallel:

```groovy
parallel {
    stage('Test UI') { ... }
    stage('Test Server') { ... }
}
```

## Performance Optimization

### 1. Cache Dependencies

```powershell
# npm cache
npm ci --prefer-offline

# Go modules
go env -w GO111MODULE=on
```

### 2. Use SSD

- Place Jenkins WORKSPACE on SSD for faster I/O

### 3. Allocate Resources

- Give Jenkins agent sufficient CPU and RAM
- Docker build: Minimum 4GB RAM, 2 CPU cores
- Go tests: Minimum 2GB RAM

### 4. Cleanup

Pipeline automatically cleans workspace with `deleteDir()`. To keep workspace:

```groovy
deleteDir()  // Remove this line
```

## Security

### 1. Credentials Management

Store sensitive data in Jenkins Credentials Store:

1. **Manage Jenkins** → **Manage Credentials**
2. Add credentials for:
   - GitHub PAT (Personal Access Token)
   - Docker Registry credentials
   - Database credentials

### 2. Secure Environment Variables

```groovy
environment {
    REGISTRY_CREDS = credentials('docker-registry-creds')
    GITHUB_TOKEN = credentials('github-token')
}
```

### 3. GitHub Webhook Security

- Use HTTPS for webhook URL
- Configure webhook secret in GitHub
- Jenkins validates webhook signature

## Support

For issues or questions:
1. Check Jenkins logs: **Manage Jenkins** → **System Logs**
2. Review pipeline console output
3. Check tool installation on agent
4. Verify network connectivity to GitHub

## References

- **Jenkins Documentation**: https://www.jenkins.io/doc/
- **Jenkins Pipeline Syntax**: https://www.jenkins.io/doc/book/pipeline/
- **Jenkinsfile Documentation**: https://www.jenkins.io/doc/book/pipeline/jenkinsfile/
- **Windows Agents Guide**: https://www.jenkins.io/doc/book/using-jenkins-agents/windows/
