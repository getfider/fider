import 'reflect-metadata';
import getDecorators from 'inversify-inject-decorators';
import { Container, injectable } from 'inversify';
import {
  IdeaService,
  HttpIdeaService,
  TenantService,
  HttpTenantService,
  UserService,
  HttpUserService,
  TagService,
  HttpTagService,
  Cache,
  BrowserCache,
} from '@fider/services';

const container = new Container();

const {
    lazyInject,
    lazyInjectNamed,
    lazyInjectTagged,
    lazyMultiInject
} = getDecorators(container);

const injectables = {
    Cache: Symbol('Cache'),
    IdeaService: Symbol('IdeaService'),
    UserService: Symbol('UserService'),
    TenantService: Symbol('TenantService'),
    TagService: Symbol('TagService'),
};

export {
    container,
    injectables,
    injectable,
    lazyInject as inject
};

container.bind<Cache>(injectables.Cache).toConstantValue(new BrowserCache(window));
container.bind<IdeaService>(injectables.IdeaService).to(HttpIdeaService);
container.bind<TenantService>(injectables.TenantService).to(HttpTenantService);
container.bind<UserService>(injectables.UserService).to(HttpUserService);
container.bind<TagService>(injectables.TagService).to(HttpTagService);
