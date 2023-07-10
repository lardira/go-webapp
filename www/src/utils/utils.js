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
}
