import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Refund Policy - MiCartaPro',
  description: 'Refund policy for MiCartaPro. Learn about our refund and cancellation policies.',
}

export default function RefundPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link href="/" className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </Link>
            <Link 
              href="/"
              className="text-gray-600 hover:text-blue-600 transition-colors"
            >
              Back to home
            </Link>
          </div>
        </div>
      </nav>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Refund Policy
          </h1>
          
          <p className="text-gray-600 mb-8">
            <strong>Last updated:</strong> {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Overview</h2>
              <p>
                This Refund Policy ("Policy") describes the terms and conditions under which MiCartaPro ("we", "us", or "our") provides refunds for our digital menu services. By purchasing our services, you agree to this Refund Policy.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Service Description</h2>
              <p>
                MiCartaPro provides digital menu services including:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Customized digital menu design and development</li>
                <li>Exclusive QR code generation</li>
                <li>Shopping cart integration</li>
                <li>WhatsApp integration for order reception</li>
                <li>Responsive design for all devices</li>
                <li>First year free (promotional offer)</li>
                <li>Monthly subscription renewal from the second year ($10 USD/month)</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Initial Setup Fee</h2>
              <p>
                Our services require an initial setup fee of $150 USD (or equivalent in your local currency). This fee covers:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Custom design and development of your digital menu</li>
                <li>Logo customization</li>
                <li>QR code generation</li>
                <li>Initial configuration and setup</li>
                <li>First year of service (free promotional period)</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Refund Eligibility</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.1. Refund Period</h3>
              <p>
                You may request a full refund of the initial setup fee within <strong>14 days</strong> of the purchase date, provided that:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>The service has not been fully delivered or activated</li>
                <li>You have not used the service extensively</li>
                <li>The refund request is made in good faith</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.2. Non-Refundable Situations</h3>
              <p>
                Refunds will NOT be provided in the following circumstances:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>After 14 days from the purchase date</li>
                <li>If the service has been fully delivered and activated</li>
                <li>If you have received and used the custom design, QR code, or other deliverables</li>
                <li>For monthly subscription fees after the service period has commenced</li>
                <li>If the refund request is due to a change of mind after service delivery</li>
                <li>For any additional services or customizations requested beyond the initial scope</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Monthly Subscription Refunds</h2>
              <p>
                Monthly subscription fees ($10 USD/month starting from the second year) are charged in advance and are generally non-refundable. However, we may consider prorated refunds in exceptional circumstances, such as:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Service unavailability for extended periods (more than 7 consecutive days) due to our fault</li>
                <li>Technical issues that prevent you from using the service and that we are unable to resolve</li>
                <li>Duplicate charges due to billing errors on our part</li>
              </ul>
              <p className="mt-4">
                Refund requests for monthly subscriptions must be submitted within 7 days of the charge date.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. How to Request a Refund</h2>
              <p>
                To request a refund, please contact us through:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> support@micartapro.com</li>
              </ul>
              <p className="mt-4">
                Your refund request must include:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Your account information or order number</li>
                <li>The reason for the refund request</li>
                <li>The date of purchase</li>
                <li>Any relevant documentation or screenshots</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Refund Processing</h2>
              <p>
                Once we receive and approve your refund request:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>We will process the refund within <strong>5-10 business days</strong></li>
                <li>The refund will be issued to the original payment method used for the purchase</li>
                <li>You will receive a confirmation email once the refund has been processed</li>
                <li>The time it takes for the refund to appear in your account depends on your payment provider (typically 3-5 business days for credit cards, 5-10 business days for bank transfers)</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Cancellation Policy</h2>
              <p>
                You may cancel your subscription at any time. Cancellation will take effect at the end of your current billing period. You will continue to have access to the service until the end of the paid period.
              </p>
              <p>
                To cancel your subscription, contact us through WhatsApp or email. No refunds will be provided for the remaining period of an active subscription unless otherwise specified in this Policy.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Chargebacks</h2>
              <p>
                If you file a chargeback or dispute with your payment provider, we reserve the right to:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Immediately suspend or terminate your account</li>
                <li>Dispute the chargeback with your payment provider</li>
                <li>Provide evidence of service delivery and this Refund Policy</li>
              </ul>
              <p className="mt-4">
                We encourage you to contact us directly to resolve any issues before initiating a chargeback, as we are committed to resolving disputes amicably.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Service Modifications</h2>
              <p>
                If we make significant changes to our services that materially affect your use of the service, you may be eligible for a prorated refund. We will notify you of any such changes in advance and provide options for cancellation and refund if applicable.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Currency and Exchange Rates</h2>
              <p>
                All refunds will be processed in the same currency as the original payment. If currency conversion is required, the exchange rate at the time of the refund will apply. We are not responsible for any fees charged by your bank or payment provider for currency conversion.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Disputes and Resolution</h2>
              <p>
                If you are not satisfied with our refund decision, you may:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Contact us to discuss your concerns further</li>
                <li>Request a review of your case by our management team</li>
                <li>Seek resolution through alternative dispute resolution mechanisms as provided by applicable law</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">13. Changes to This Policy</h2>
              <p>
                We reserve the right to modify this Refund Policy at any time. Changes will be effective immediately upon posting on this page. We will notify you of material changes via email or through our Services.
              </p>
              <p>
                Your continued use of our Services after changes to this Policy constitutes acceptance of the modified Policy.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">14. Contact Information</h2>
              <p>
                For questions about this Refund Policy or to request a refund, please contact us:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> support@micartapro.com</li>
                <li><strong>Business Hours:</strong> Monday to Friday, 9:00 AM - 5:00 PM (Chile Time)</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">15. Governing Law</h2>
              <p>
                This Refund Policy is governed by and construed in accordance with applicable consumer protection laws and regulations. Any disputes arising from this Policy will be resolved in accordance with the laws of your jurisdiction and applicable international consumer protection standards.
              </p>
            </section>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-200 flex flex-col sm:flex-row justify-between items-center gap-4">
            <Link 
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              ← Back to home
            </Link>
            <Link 
              href="/terms-en"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              View Terms and Conditions →
            </Link>
          </div>
        </div>
      </div>
    </div>
  )
}

