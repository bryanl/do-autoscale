import Ember from 'ember';

export default Ember.Route.extend({
  model() {
    return Ember.RSVP.hash({
      groupConfig: this.store.queryRecord('group-config', {})
    });
  },
  flashMessages: Ember.inject.service(),
  actions: {
    createGroup(options) {
      var group = this.store.createRecord('group', options);
      return group.save().then(() => {
        this.transitionTo('groups');
      }).catch((reason) => {
        const flashMessages = this.get('flashMessages');

        console.log(reason);
        flashMessages.danger("Unable to create group!");
      });
    }
  }

});
