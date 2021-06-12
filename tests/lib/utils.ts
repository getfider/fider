import "reflect-metadata"

export const delay = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms))
}

export function findBy(selector: string) {
  return (target: any, propertyKey: string) => {
    const type = Reflect.getMetadata("design:type", target, propertyKey)
    Object.defineProperty(target, propertyKey, {
      configurable: true,
      enumerable: true,
      get() {
        const tab = this.tab
        return new type(tab, selector)
      },
    })
  }
}
