import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    return Ember.RSVP.hash({
      groupConfig: this.store.queryRecord('group-config', {})
    });
  },

  actions: {
    createGroup(options) {
     var group = this.store.createRecord('group', options);
      return group.save();
    }
  }

});
