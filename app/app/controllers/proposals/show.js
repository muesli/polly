import Ember from 'ember';

const { inject: { service }, Component } = Ember;

export default Ember.Controller.extend({
  currentUser: service('current-user'),

  actions: {
    moderate(id) {
        this.store.findRecord('proposal', id).then(function(proposal) {
            proposal.set('moderated', true)
            proposal.save();
        });
    },
    vote(id) {
        alert(id);
    }
  }
});
