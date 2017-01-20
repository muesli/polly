import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    model(params) {
       return Ember.RSVP.hash({
         proposal: this.store.findRecord('proposal', params.proposal_id)
       });
     },

     setupController(controller, models) {
       controller.set('proposal', models.proposal);
     }
});
