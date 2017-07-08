import * as React from 'react';
import * as moment from 'moment';

interface MomentText {
  date: Date | string;
}

export const Moment = (props: MomentText) => {
  if (!props.date) {
    return <span></span>;
  }

  const m = moment(props.date);
  const diff = m.diff(new Date(), 'years');
  const display = (diff !== 0) ? m.format('Do MMM YYYY') : m.fromNow();

  return <span className="date" title={m.format('Do MMM YYYY, hh:mm')}>
    { display }
  </span>;
};
