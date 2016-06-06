import Ember from 'ember';

export default Ember.Route.extend({
  model(params) {
    return Ember.RSVP.hash({
      group: this.store.findRecord('group', params.group_id)
    })
   },

   currentModel: function () { return this.modelFor(this.routeName); },

   actions: {
  }
});
