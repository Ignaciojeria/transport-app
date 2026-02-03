package events

const (
	EventUserMenusInsertedWebhook         = "user.menus.inserted.webhook"
	EventMenuInteractionRequested         = "menu.interaction.requested"
	EventMenuCreateRequested              = "menu.create.requested"
	EventCreemSubscriptionTrialingWebhook = "creem.subscription.trialing.webhook"
	EventCreemCheckoutCompletedWebhook    = "creem.checkout.completed.webhook"
	EventCreemSubscriptionActiveWebhook   = "creem.subscription.active.webhook"
	EventCreemSubscriptionPaidWebhook     = "creem.subscription.paid.webhook"
	EventCreemSubscriptionCanceledWebhook = "creem.subscription.canceled.webhook"
	EventCreemSubscriptionExpiredWebhook  = "creem.subscription.expired.webhook"
	EventCreemSubscriptionUpdateWebhook   = "creem.subscription.update.webhook"
	EventCreemSubscriptionPausedWebhook   = "creem.subscription.paused.webhook"
	EventCreemRefundCreatedWebhook        = "creem.refund.created.webhook"
	EventCreemDisputeCreatedWebhook       = "creem.dispute.created.webhook"
	EventCreateOrderRequested             = "create.order.requested"
	EventOrderStartedPreparation          = "order.started.preparation"
	EventOrderItemReady                   = "order.item.ready"
	EventOrderDispatched                  = "order.dispatched"
	EventOrderCancelled                   = "order.cancelled"
	EventImageGenerationRequested         = "image.generation.requested"
	EventImageEditionRequested           = "image.edition.requested"
)
