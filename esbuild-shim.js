const global = (1, eval)('this');
global.global = global;
global.globalThis = global;
global.frames = global;
global.self = global;

const window = {}
const document = {
  documentElement: {},
  getElementById: () => undefined
}
const navigator = {}

global.navigator = navigator;
global.window = window;
global.window.document = document;
global.window.location = { href: '' };
global.document = document;