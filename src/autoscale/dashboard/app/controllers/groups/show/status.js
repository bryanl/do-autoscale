import Ember from 'ember';

export default Ember.Controller.extend({

  data: {
    x: 'x',
    //        xFormat: '%Y%m%d', // 'xFormat' can be used as custom format of 'x'
    columns: [
      ['x', '2013-01-01', '2013-01-02', '2013-01-03', '2013-01-04', '2013-01-05', '2013-01-06'],
      //            ['x', '20130101', '20130102', '20130103', '20130104', '20130105', '20130106'],
      ['data1', 30, 200, 100, 400, 150, 250],
      ['data2', 130, 340, 200, 500, 250, 350]
    ]
  },

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
