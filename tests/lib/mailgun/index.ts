require('dotenv').config();
import axios from 'axios';
import * as https from 'https';
import { delay } from '../';

const httpGet = (url: string): Promise<any> => {
  return new Promise((resolve, reject) => {
    https.get(url, (res: https.IncomingMessage) => {
      // TODO: implement this, add authentication and remove axios
    });
  });
};

const mailgunOptions = {
  auth: {
    username: 'api',
    password: process.env.MAILGUN_API!
  }
};

export const mailgun = {
  getLinkFromLastEmailTo: async (to: string): Promise<string> => {
    let messageUrl = '';
    let count = 0;

    do {
      count++;
      const url = `https://api.mailgun.net/v3/${process.env.MAILGUN_DOMAIN}/events?to=${to}&event=accepted&limit=1&ascending=no`;
      const events = await axios.get(url, mailgunOptions);
      if (events.data.items.length > 0 && events.data.items[0].message.headers.to === to) {
        messageUrl = events.data.items[0].storage.url;
      } else {
        await delay(5000);
      }
    } while (!messageUrl && count < 6);

    if (count === 6) {
      throw new Error(`Message not found for ${to}.`);
    }

    const message = await axios.get(messageUrl, mailgunOptions);

    const matches = /<a\s+(?:[^>]*?\s+)?href=(["'])(.*?)\1/.exec(message.data['body-html']);
    if (matches) {
      return matches[2];
    }

    return '';
  }
};
