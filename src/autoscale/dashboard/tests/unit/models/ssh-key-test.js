import { moduleForModel, test } from 'ember-qunit';

moduleForModel('ssh-key', 'Unit | Model | ssh key', {
  // Specify the other units that are required for this test.
  needs: []
});

test('it exists', function(assert) {
  let model = this.subject();
  // let store = this.store();
  assert.ok(!!model);
});
