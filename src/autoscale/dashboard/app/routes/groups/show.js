import Ember from 'ember';

export default Ember.Route.extend({
   beforeModel: function() {
    this.transitionTo('groups.show.status');
  },

  model(params) {
      return this.store.findRecord('group', params.group_id);
   },

   currentModel: function () { return this.modelFor(this.routeName); },

   actions: {
  }
});
