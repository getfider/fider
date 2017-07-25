const config = {
  baseUrl: process.env.BASE_URL || 'http://demo.dev.fider.io:3000',
  users: {
    darthvader: {
      email: 'darthvader.fider@gmail.com',
      password: process.env.DARTHVADER_PASSWORD!,
    },
  },
};

export default config;
