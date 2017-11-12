import * as React from 'react';
import * as moment from 'moment';

interface MomentText {
  date: Date | string;
}

export const Moment = (props: MomentText) => {
  if (!props.date) {
    return <span />;
  }

  const m = moment(props.date);
  const diff = m.diff(new Date(), 'years');
  const display = (diff !== 0) ? m.format('MMMM D, YYYY') : m.fromNow();

  return (
    <span className="date" title={m.format('MMMM D, YYYY')}>
      {display}
    </span>
  );
};
