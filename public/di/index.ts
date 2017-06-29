import 'reflect-metadata';
import getDecorators from 'inversify-inject-decorators';
import { Container, injectable } from 'inversify';

const container = new Container();

const {
    lazyInject,
    lazyInjectNamed,
    lazyInjectTagged,
    lazyMultiInject
} = getDecorators(container);

const injectables = {
    Session: Symbol('Session'),
    IdeaService: Symbol('IdeaService'),
    TenantService: Symbol('TenantService'),
};

export {
    container,
    injectables,
    injectable,
    lazyInject as inject
};
