import attr from 'ember-data/attr';
import Fragment from 'model-fragments/fragment';

export default Fragment.extend({
  min_size: attr('number'),
  max_size: attr('number'),
  scale_up_value: attr(),
  scale_up_by: attr(),
  scale_down_value: attr(),
  scale_down_by: attr(),
  warm_up_duration: attr()
});
