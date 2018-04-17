import { toast } from "react-toastify";

export const success = (message: string | JSX.Element) => {
  toast.success(message);
};

export const error = (message: string | JSX.Element) => {
  toast.error(message);
};
