require('dotenv').config();
import axios from 'axios';

export const mailgun = {
  getLinkFromLastEmailTo: async (to: string): Promise<string> => {
    const url = `https://api.mailgun.net/v3/${process.env.MAILGUN_DOMAIN}/events?to=${to}&event=accepted&limit=1&ascending=no`;
    let response = await axios.get(url, {
      auth: {
        username: 'api',
        password: process.env.MAILGUN_API!
      }
    });

    response = await axios.get(response.data.items[0].storage.url, {
      auth: {
        username: 'api',
        password: process.env.MAILGUN_API!
      }
    });

    const matches = /<a\s+(?:[^>]*?\s+)?href=(["'])(.*?)\1/.exec(response.data['body-html']);
    if (matches) {
      return matches[2];
    }

    return '';
  }
};
