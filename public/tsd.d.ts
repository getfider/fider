interface Window {
  ga?: (cmd: string, evt: string, args?: any) => void;
}

declare var require: (id: string) => any;
