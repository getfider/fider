import * as React from 'react';

interface MomentText {
  date: Date | string;
}

const templates: { [key: string]: string } = {
  seconds: 'less than a minute',
  minute: 'about a minute',
  minutes: '%d minutes',
  hour: 'about an hour',
  hours: 'about %d hours',
  day: 'a day',
  days: '%d days',
  month: 'about a month',
  months: '%d months',
  year: 'about a year',
  years: '%d years'
};

const template = (t: string, n: number): string => {
  return templates[t] && templates[t].replace(/%d/i, Math.abs(Math.round(n)).toString());
};

const monthNames = [
  'January', 'February', 'March',
  'April', 'May', 'June', 'July',
  'August', 'September', 'October',
  'November', 'December'
];

function format(date: Date) {

  const day = date.getDate();
  const monthIndex = date.getMonth();
  const year = date.getFullYear();

  return `${monthNames[date.getMonth()]} ${date.getDate()}, ${date.getFullYear()} · ${date.getHours()}:${date.getMinutes()}`;
}

function timeSince(now: Date, date: Date) {

  const seconds = (now.getTime() - date.getTime()) / 1000;
  const minutes = seconds / 60;
  const hours = minutes / 60;
  const days = hours / 24;
  const years = days / 365;

  return (
    seconds < 45 && template('seconds', seconds) ||
    seconds < 90 && template('minute', 1) ||
    minutes < 45 && template('minutes', minutes) ||
    minutes < 90 && template('hour', 1) ||
    hours < 24 && template('hours', hours) ||
    hours < 42 && template('day', 1) ||
    days < 30 && template('days', days) ||
    days < 45 && template('month', 1) ||
    days < 365 && template('months', days / 30) ||
    years < 1.5 && template('year', 1) ||
    template('years', years)
  ) + ' ago';
}

export const Moment = (props: MomentText) => {
  if (!props.date) {
    return <span />;
  }

  const now = new Date();
  const date = props.date instanceof Date ? props.date : new Date(props.date);

  const diff = (now.getTime() - date.getTime()) / (60 * 60 * 24 * 1000);

  // const m = moment(props.date);
  // const diff = m.diff(new Date(), 'years');
  // const display = (diff !== 0) ? m.format('MMMM D, YYYY') : m.fromNow();
  const display = (diff >= 365) ? format(date) : timeSince(now, date);

  return (
    <span className="date" title={format(date)}>
      {display}
    </span>
  );
};
