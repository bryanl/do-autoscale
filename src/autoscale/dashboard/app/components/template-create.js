import Ember from 'ember';

export default Ember.Component.extend({
  init() {
    this._super(...arguments);
    console.log("starting up");
    console.log(this.userConfig);
  },

  // should the create form be displayed?
  showCreateForm: false,

  // form items
  name: "template-1",
  region: "tor1",
  size: "4gb",
  image: "ubuntu-14-04-x64",
  userData: null,
  sshKeys: [{id: 104064}],

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

      console.log(createRequest);

      const promise = this.get("onCreate")(createRequest);
      promise.then(() => {
        this.set("show", false);
      });
    }
  }
});
