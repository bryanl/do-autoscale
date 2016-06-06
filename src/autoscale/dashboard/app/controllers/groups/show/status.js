import Ember from 'ember';

export default Ember.Controller.extend({

  newData: Ember.computed('model', function () {
    const model = this.get('model');

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
    }

    return obj
  }),

  axis: {
    x: {
      type: 'timeseries',
      tick: {
        format: '%Y-%m-%d %I:%M',
        culling: true,
      }
    }
  }
});
