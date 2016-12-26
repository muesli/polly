import Ember from 'ember';
import AuthenticatedRouteMixin from 'ember-simple-auth/mixins/authenticated-route-mixin';

export default Ember.Route.extend(AuthenticatedRouteMixin, {
    model() {
        return Ember.RSVP.hash({
            users: this.get('store').findAll('user')
        });
    },

    setupController(controller, models) {
        controller.set('users', models.users);
    }
});
