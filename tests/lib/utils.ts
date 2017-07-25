import 'reflect-metadata';
import { ThenableWebDriver, WebElementPromise } from 'selenium-webdriver';

export const delay = (ms: number) => {
    return new Promise((resolve) => setTimeout(resolve, ms));
};

export const timeout = (ms: number) => {
    return new Promise((_, reject) => setTimeout(() => reject(new Error('timeout')), ms));
};

export function findBy(selector: string, t?: new(element: WebElementPromise, selector: string) => any) {
  return (target: any, propertyKey: string) => {
    const type = Reflect.getMetadata('design:type', target, propertyKey);
    Object.defineProperty(target, propertyKey, {
        configurable: true,
        enumerable: true,
        get() {
          const promise = ((this as any).browser as ThenableWebDriver).findElement(selector);
          return new type(promise, selector);
        },
    });
  };
}

export function findMultipleBy(selector: string, t?: new(element: WebElementPromise, selector: string) => any) {
  return (target: any, propertyKey: string) => {
    const type = Reflect.getMetadata('design:type', target, propertyKey);
    Object.defineProperty(target, propertyKey, {
        configurable: true,
        enumerable: true,
        get() {
          const promise = ((this as any).browser as ThenableWebDriver).findElements(selector);
          return new type(promise, selector);
        },
    });
  };
}
