import Ember from 'ember';

export default Ember.Route.extend({

  model() {
    return Ember.RSVP.hash({
      templates: this.store.findAll('template'),
      userConfig: this.store.queryRecord('user_config', {})
    });
  },


  createTemplate(options) {
    console.log("creating new template (up top)");
    console.log(options);

    var template = this.store.createRecord('template', options);
    return template.Save();
  },

  actions: {
    createTemplate: function(options) {
      console.log("creating new template");
      console.log(options);

      var template = this.store.createRecord('template', options);
      console.log(template);
      return template.save();
    }
  }

});
