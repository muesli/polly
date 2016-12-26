import { moduleForModel, test } from 'ember-qunit';

moduleForModel('proposal', 'Unit | Model | proposal', {
  // Specify the other units that are required for this test.
  needs: []
});

test('it exists', function(assert) {
  let model = this.subject();
  // let store = this.store();
  assert.ok(!!model);
});
