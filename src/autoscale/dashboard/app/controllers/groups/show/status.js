import Ember from 'ember';

export default Ember.Controller.extend({

  afterModel: function(model) {
    model.group.reload();
  },

  resourceStatus: Ember.computed('model', function () {
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

    var now = new Date();
    var lines = [
      { value: this.relativeDate(now, 270) },
      { value: this.relativeDate(now, 180) },
      { value: this.relativeDate(now, 90) },
      { value: this.relativeDate(now, 0) }];

    var obj = {
      data: {
        type: 'line',
        json: model.get('scaleHistory').serialize(),
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
            values: yValues,
            count: yValues.length
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


  metricStatus: Ember.computed('model.group.timeseriesValues', function () {
    var values = this.get('model.group.timeseriesValues');

    var json = [];
    values.forEach((x) => {
      json.push(x.serialize());
    });

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
            // values: yValues,
            count: 4,
          }
        }
      },

      grid: {
        x: {
          // lines: lines,
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
