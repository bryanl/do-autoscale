
import Ember from 'ember';

export default Ember.Route.extend({
  currentModel: function () { return this.modelFor(this.routeName); },
  flashMessages: Ember.inject.service(),
  actions: {
    submit() {
      const flashMessages = this.get('flashMessages');

      this.currentModel.group.save().then(function(){
        flashMessages.success("Group updated!");
      }).catch(function(reason){
        console.log(reason);
        flashMessages.danger("Unable to update group!");
      });
    }
  }
});
