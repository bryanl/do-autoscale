import Ember from 'ember';

export default Ember.Component.extend({
  // should the create form be displayed?
  showCreateForm: false,

  actions: {
    addTemplate: function() {
      this.set("showCreateForm", true);
    },
    createTemplate(options) {
      console.log("in template list");
      console.log(options);
      this.sendAction(this.get("onTemplateCreate"), options);
    }
  }
});
