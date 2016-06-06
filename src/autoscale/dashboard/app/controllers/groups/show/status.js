import Ember from 'ember';

export default Ember.Controller.extend({

  data: Ember.computed('model', function () {
    const model = this.get('model').group;

    var obj = {
      type: 'line',
      json: model.get('scaleHistory').serialize(),
      keys: {
        x: 'createdAt',
        xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
        value: ['total'],
      },
      x: 'x',
      xFormat: '%Y-%m-%dT%H:%M:%S.%LZ'
    };

    return obj;
  }),

  axis: Ember.computed('model', function () {
    const model = this.get('model').group;

    var maxY = 0;
    var histories = model.get('scaleHistory').serialize();
    for (let history of histories) {
      if (history.total > maxY) {
        maxY = history.total;
      }
    }

    var yValues = [];
    maxY = maxY + (maxY * 0.10);

    for (var i = 0; i < 4; i++) {
      yValues[i] = Math.ceil(maxY) / 3 * i;
    }
    var obj = {
      x: {
        type: 'timeseries',
        tick: {
          format: '%Y-%m-%d %I:%M',
          count: 4
        }
      },
      y: {
        tick: {
          values: yValues,
          count: yValues.length
        }
      }
    };

    return obj;
  }),

  grid: Ember.computed('model', function() {
    var now = new Date();
    var lines = [
      {value: this.relativeDate(now, 270)},
      {value: this.relativeDate(now, 180)},
      {value: this.relativeDate(now, 90)},
      {value: this.relativeDate(now, 0)}];

    var obj = {
      x: {
        lines: lines,
      },
      y: {
        show: true
      }
    };

    return obj;
  }),

  relativeDate: function(now, relativeTime) {
    // var in_utc = new Date(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate(),  now.getUTCHours(), now.getUTCMinutes(), now.getUTCSeconds());
    return (new Date(now - relativeTime * 60 * 1000)).toISOString();
  }
});
