import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    return Ember.RSVP.hash({
      groups: this.store.findAll('group'),
    });
  },

  actions: {
    addGroup() {
      this.transitionTo('groups.new');
    }
  }
});
