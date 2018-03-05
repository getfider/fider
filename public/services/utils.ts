export const delay = (ms: number) => {
  return new Promise((resolve) => setTimeout(resolve, ms));
};

export const classSet = (input?: any): string => {
  let classes = '';
  if (input) {
    for (const key in input) {
      if (key && !!input[key]) {
        classes += ` ${key}`;
      }
    }
    return classes.trimLeft();
  }
  return '';
};
