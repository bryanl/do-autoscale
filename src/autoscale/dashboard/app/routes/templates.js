import Ember from 'ember';

export default Ember.Route.extend({

  model() {
    return [
      { id: "1", name: "template-1", region: "nyc1", ssh_keys: ["1", "2"], user_data: "#data" }
    ];
  }
});
