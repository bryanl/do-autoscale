import Model from 'ember-data/model';
import attr from 'ember-data/attr';

export default Model.extend({
  name: attr(),
  region: attr(),
  size: attr(),
  image: attr(),
  sshKeys: attr(),
  userData: attr()
});
