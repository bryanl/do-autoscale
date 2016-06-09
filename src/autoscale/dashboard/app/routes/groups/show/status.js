import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    var model = this.modelFor("groups.show");
    model.reload();
    return model;
  },

});
