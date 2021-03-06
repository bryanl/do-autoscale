import Ember from 'ember';

export default Ember.Controller.extend({
  resourceStatus: Ember.computed('model', function () {
    const group = this.get('model');

    var maxY = 0;
    var histories = group.get('scaleHistory').serialize();
    for (let history of histories) {
      if (history.total > maxY) {
        maxY = history.total;
      }
    }

    var now = new Date();
    var lines = [
      { value: this.relativeDate(now, 270) },
      { value: this.relativeDate(now, 180) },
      { value: this.relativeDate(now, 90) },
      { value: this.relativeDate(now, 0) }];

    var obj = {
      data: {
        type: 'line',
        json: group.get('scaleHistory').serialize(),
        keys: {
          x: 'createdAt',
          xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
          value: ['total'],
        },
        x: 'x',
        xFormat: '%Y-%m-%dT%H:%M:%S.%LZ'
      },

      axis: {
        x: {
          type: 'timeseries',
          tick: {
            format: '%Y-%m-%d %I:%M',
            count: 4
          }
        },
        y: {
          tick: {
            format: function (d) { return d.toFixed(2).replace(/\.?0*$/, ''); }
          }
        }
      },

      grid: {
        x: {
          lines: lines,
        },
        y: {
          show: true
        }
      },

      point: {
        show: false,
      }
    };

    return obj;
  }),


  metricStatus: Ember.computed('model.timeseriesValues', function () {
    var values = this.get('model.timeseriesValues');

    var json = [];
    values.forEach((x) => {
      json.push(x.serialize());
    });

    var now = new Date();
    var lines = [
      { value: this.relativeDate(now, 270) },
      { value: this.relativeDate(now, 180) },
      { value: this.relativeDate(now, 90) },
      { value: this.relativeDate(now, 0) }];

    var obj = {
      data: {
        type: 'line',
        json: json,
        keys: {
          x: 'timestamp',
          xFormat: '%Y-%m-%dT%H:%M:%S.%LZ',
          value: ['value'],
        },
        x: 'x',
        xFormat: '%Y-%m-%dT%H:%M:%S.%LZ'
      },

      axis: {
        x: {
          type: 'timeseries',
          tick: {
            format: '%Y-%m-%d %I:%M',
            count: 4
          }
        },
        y: {
          tick: {
            count: 4,
            format: function (d) { return d.toFixed(2).replace(/\.?0*$/, ''); }

          }
        }
      },

      grid: {
        x: {
          lines: lines,
        },
        y: {
          show: true
        }
      },

      point: {
        show: false,
      }

    };

    return obj;
  }),

  relativeDate: function (now, relativeTime) {
    return (new Date(now - relativeTime * 60 * 1000)).toISOString();
  }
});
