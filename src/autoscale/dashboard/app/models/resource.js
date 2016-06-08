import attr from 'ember-data/attr';
import Fragment from 'model-fragments/fragment';

export default Fragment.extend({
  address: attr('string'),
  name: attr('string'),
  createdAt: attr('date')
});
