import Ember from 'ember';

export default Ember.Component.extend({
  init() {
    this._super(...arguments);
    this.set('sshKeys', this.userConfig.get('keys'));
  },

  // should the create form be displayed?
  showCreateForm: false,

  // form items
  name: "template-1",
  size: "4gb",
  image: "ubuntu-14-04-x64",
  userData: null,

  actions: {
    submit: function() {
      var createRequest = {
        name: this.name,
        region: this.region,
        size: this.size,
        image: this.image,
        sshKeys: this.sshKeys,
        userData: this.userData
      };

      const promise = this.get("onCreate")(createRequest);
      promise.then(() => {
        this.set("show", false);
      });
    }
  }
});
