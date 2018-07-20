import { Fider } from '.';

export const delay = (ms: number) => {
  return new Promise(resolve => setTimeout(resolve, ms));
};

export const classSet = (input?: any): string => {
  let classes = "";
  if (input) {
    for (const key in input) {
      if (key && !!input[key]) {
        classes += ` ${key}`;
      }
    }
    return classes.trim();
  }
  return "";
};

const monthNames = [
  "January",
  "February",
  "March",
  "April",
  "May",
  "June",
  "July",
  "August",
  "September",
  "October",
  "November",
  "December"
];

const twoDigits = (value: number): string => {
  return value < 9 ? `0${value}` : value.toString();
};

export const formatDate = (input: Date | string): string => {
  const date = input instanceof Date ? input : new Date(input);

  const day = date.getDate();
  const monthIndex = date.getMonth();
  const year = date.getFullYear();
  const hours = twoDigits(date.getHours());
  const minutes = twoDigits(date.getMinutes());

  return `${monthNames[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()} · ${hours}:${minutes}`;
};

const templates: { [key: string]: string } = {
  seconds: "less than a minute",
  minute: "about a minute",
  minutes: "%d minutes",
  hour: "about an hour",
  hours: "about %d hours",
  day: "a day",
  days: "%d days",
  month: "about a month",
  months: "%d months",
  year: "about a year",
  years: "%d years"
};

const template = (t: string, n: number): string => {
  return templates[t] && templates[t].replace(/%d/i, Math.abs(Math.round(n)).toString());
};

export const timeSince = (now: Date, date: Date): string => {
  const seconds = (now.getTime() - date.getTime()) / 1000;
  const minutes = seconds / 60;
  const hours = minutes / 60;
  const days = hours / 24;
  const years = days / 365;

  return (
    ((seconds < 45 && template("seconds", seconds)) ||
      (seconds < 90 && template("minute", 1)) ||
      (minutes < 45 && template("minutes", minutes)) ||
      (minutes < 90 && template("hour", 1)) ||
      (hours < 24 && template("hours", hours)) ||
      (hours < 42 && template("day", 1)) ||
      (days < 30 && template("days", days)) ||
      (days < 45 && template("month", 1)) ||
      (days < 365 && template("months", days / 30)) ||
      (years < 1.5 && template("year", 1)) ||
      template("years", years)) + " ago"
  );
};

export const fileToBase64 = async (file: File): Promise<string> => {
  return new Promise<string>((resolve, reject) => {
    const reader = new FileReader();
    reader.addEventListener(
      "load",
      () => {
        const parts = (reader.result as string).split("base64,");
        resolve(parts[1]);
      },
      false
    );

    reader.addEventListener(
      "error",
      () => {
        reject(reader.error);
      },
      false
    );

    reader.readAsDataURL(file);
  });
};

export const isCookieEnabled = (): boolean => {
  try {
    document.cookie = "cookietest=1";
    const ret = document.cookie.indexOf("cookietest=") !== -1;
    document.cookie = "cookietest=1; expires=Thu, 01-Jan-1970 00:00:01 GMT";
    return ret;
  } catch (e) {
    return false;
  }
};

export const uploadedImageUrl = (id: number, size: number): string | undefined => {
  if (id > 0) {
    return `${Fider.settings.assetsURL}/images/${size}/${id}`;
  }
  return undefined;
};
