import Model from 'ember-data/model';
import attr from 'ember-data/attr';
import { fragment, fragmentArray } from 'model-fragments/attributes';

export default Model.extend({
  name: attr(),
  baseName: attr(),
  templateID: attr(),
  metricType: attr(),
  metric: attr(),
  policyType: attr(),
  policy: fragment('policy'),
  scaleHistory: fragmentArray('group-status'),
  timeseriesValues: fragmentArray('timeseries'),
  resources: fragmentArray('resource')
});
