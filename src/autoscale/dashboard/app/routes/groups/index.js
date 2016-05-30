import Ember from 'ember';

export default Ember.Route.extend({
  actions: {
    addGroup() {
      this.transitionTo('groups.new');
    }
  }
});
