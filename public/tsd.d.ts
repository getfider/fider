interface JQuery {
  modal(opts: any): any;
  popup(args: any): any;
  dropdown(args?: any, args2?: any): any;
}

interface Window {
  //TODO: implement correct API
  ga?: (cmd: string, evt: string) => void;
}