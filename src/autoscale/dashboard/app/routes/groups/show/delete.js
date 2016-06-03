import Ember from 'ember';

export default Ember.Route.extend({
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
