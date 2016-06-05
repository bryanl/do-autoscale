import attr from 'ember-data/attr';
import Fragment from 'model-fragments/fragment';

export default Fragment.extend({
  groupID: attr(),
  delta: attr(),
  total: attr(),
  createdAt: attr('date')
});
