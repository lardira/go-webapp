class MsgAggregator {
  constructor() {
    this.msg = '';
  }

  add(msg) {
    this.msg += msg + '\n';
  }

  empty() {
    return !this.msg;
  }

  aggregate() {
    return this.msg;
  }
}

export default class Utils {
  static newAggregator() {
    return new MsgAggregator();
  }

  static async fetchApi(method, url, body = undefined, apiKey = undefined) {
    const requestOptions = {
      method: method.toUpperCase(),
      headers: { 'Content-Type': 'application/json' },
      body: body && body instanceof String ? body : JSON.stringify(body),
    };

    if (apiKey) {
      requestOptions.headers['Authorization'] = apiKey;
    }

    return await fetch(url, requestOptions)
      .then(
        (response) =>
          new Promise((resolve, reject) => {
            const contentType = response.headers.get('Content-Type');
            if (!response.ok) response.text().then(reject);
            else {
              switch (contentType) {
                case 'application/json':
                  response.json().then(resolve);
                  break;
                default:
                  response.text().then(resolve);
              }
            }
          })
      )
      .catch(alert);
  }

  static arraysEqual(a, b) {
    if (a === b) return true;
    if (a == null || b == null) return false;
    if (a.length !== b.length) return false;
  
    for (var i = 0; i < a.length; ++i) {
      if (a[i] !== b[i]) return false;
    }
    return true;
  }
}
