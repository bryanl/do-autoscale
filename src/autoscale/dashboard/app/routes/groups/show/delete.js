import Ember from 'ember';

export default Ember.Route.extend({
   currentModel: function () { return this.modelFor(this.routeName); },

   flashMessages: Ember.inject.service(),

   actions: {
    deleteGroup() {
      const flashMessages = this.get('flashMessages');

      var group = this.currentModel.group;
      group.destroyRecord()
        .then(this.transitionTo('groups'))
        .catch(()=>{
          flashMessages.danger("Unable to remove group");
        });
    }
  }
});
