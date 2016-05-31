import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
      return this.store.findRecord('group', params.group_id);
   },

   currentModel: function () { return this.modelFor(this.routeName); },

   actions: {
    deleteGroup() {
      var group = this.currentModel;
      group.destroyRecord()
        .then(this.transitionTo('groups'))
        .catch(failure);
    }
  }
});
