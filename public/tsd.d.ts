interface JQuery {
  modal(opts: any): any;
  popup(args: any): any;
  dropdown(args?: any, args2?: any): any;
}

interface Window {
  ga?: (cmd: string, evt: string, args?: any) => void;
}

declare var require: (id: string) => any;
