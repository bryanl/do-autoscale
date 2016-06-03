import Ember from 'ember';

export default Ember.Component.extend({

  // form items
  name: "group-1",
  baseName: "doas",
  metricType: "load",
  metric: {},
  policyType: "value",
  policy: {
    min_size: 1,
    max_size: 10,
    scale_up_value: 0.8,
    scale_down_value: 0.2,
    scale_up_by: 2,
    scale_down_by: 1,
    warm_up_duration: "10s"
  },

  actions: {
    onTemplateChanged() {
      console.log(this.get('template'));
      this.set('currentTemplate', this.get('template'));
    },
    onPolicyTypeChanged() {
    },
    onMetricTypeChanged() {
    },
    submit() {
      var createRequest = {
        name: this.name,
        baseName: this.baseName,
        templateID: this.template.id,
        metricType: this.metricType,
        metric: this.metric,
        policyType: this.policyType,
        policy: this.policy
      };

      console.log(createRequest);

      const promise = this.get("onCreate")(createRequest);
      promise.then(() => {
        console.log("group created");
      });
    }
  }
});
