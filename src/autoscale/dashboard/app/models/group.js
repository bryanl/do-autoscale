import Model from 'ember-data/model';
import attr from 'ember-data/attr';

export default Model.extend({
  name: attr(),
  baseName: attr(),
  templateName: attr(),
  metricType: attr(),
  metric: attr(),
  policyType: attr(),
  policy: attr()
});
