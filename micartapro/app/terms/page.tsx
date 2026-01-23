import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Terms and Conditions - MiCartaPro',
  description: 'Terms and conditions of use for MiCartaPro. Learn about the terms governing the use of our services.',
}

export default function TermsPage() {
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
              Back to Home
            </Link>
          </div>
        </div>
      </nav>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Terms and Conditions of Use
          </h1>
          
          <p className="text-gray-600 mb-8">
            <strong>Last updated:</strong> December 3, 2025
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Acceptance of Terms</h2>
              <p>
                Welcome to MiCartaPro. These Terms and Conditions of Use ("Terms") govern your access and use of our digital menu services, including our website, applications, and any related services (collectively, the "Services").
              </p>
              <p>
                By accessing or using our Services, you agree to be bound by these Terms. If you do not agree to any part of these Terms, you must not use our Services.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Description of Services</h2>
              <p>
                MiCartaPro is a self-service SaaS (Software as a Service) platform that allows restaurants to create, manage, and share digital menus with their customers. Our platform provides:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Self-service tools for customizing digital menu designs</li>
                <li>QR code generation functionality</li>
                <li>Integrated shopping cart system</li>
                <li>WhatsApp integration for order reception</li>
                <li>Responsive design for all devices</li>
                <li>Continuous platform improvements, new features, and updates</li>
              </ul>
              <p className="mt-4">
                <strong>Important:</strong> MiCartaPro is a self-service platform. We do not provide ongoing human-driven services such as custom design work, manual menu creation, or personalized consulting as part of the subscription. All plans include access to the platform and ongoing software improvements.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Eligibility</h2>
              <p>
                To use our Services, you must:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Be at least 18 years of age or have the consent of a parent or legal guardian</li>
                <li>Have the legal capacity to enter into binding contracts</li>
                <li>Not be prohibited from using our Services under applicable laws</li>
                <li>Provide accurate and complete information when registering or using our Services</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Use of Services</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.1. Permitted Use</h3>
              <p>You may use our Services solely for lawful purposes and in accordance with these Terms. You agree to:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Use the Services responsibly and ethically</li>
                <li>Respect the intellectual property rights of others</li>
                <li>Comply with all applicable laws and regulations</li>
                <li>Not interfere with the operation of the Services</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.2. Prohibited Use</h3>
              <p>The following is strictly prohibited:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Illegal or fraudulent activities</li>
                <li>Publishing misleading pricing or illegal food products</li>
                <li>Unauthorized alcohol sales where prohibited by law</li>
                <li>Attempting unauthorized access to our systems or data</li>
                <li>Transmitting viruses, malware, or harmful code</li>
                <li>Reverse engineering, decompiling, or disassembling the Services</li>
                <li>Using bots or automated scripts to access the platform</li>
                <li>Copying, modifying, or distributing the Services without permission</li>
                <li>Using the platform to send spam or unsolicited communications</li>
                <li>Violating privacy or intellectual property rights of others</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Accounts and Registration</h2>
              <p>
                Some features require account creation. You are responsible for:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Keeping login credentials confidential</li>
                <li>All activity under your account</li>
                <li>Notifying us of any unauthorized access</li>
                <li>Providing accurate, up-to-date information</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. User Content</h2>
              <p>
                You may upload menu information, images, text, and other content ("User Content"). By doing so, you grant MiCartaPro a non-exclusive, worldwide, royalty-free, transferable license to use, reproduce, modify, and distribute such content solely for operating and improving the Service.
              </p>
              <p>
                You are responsible for ensuring your content:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Does not infringe third-party rights</li>
                <li>Is not defamatory, obscene, or illegal</li>
                <li>Complies with privacy, labeling, and food safety laws</li>
                <li>Does not violate regulations related to alcohol or consumer protection</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Intellectual Property</h2>
              <p>
                All rights in the Service, including software, design, logos, text, graphics, and assets, belong to MiCartaPro or its licensors. You receive only a limited right to use the Services under these Terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Pricing and Payments</h2>
              <ul className="list-disc pl-6 space-y-2">
                <li>Subscription pricing is shown at checkout and may change with notice</li>
                <li>No setup fees or long-term commitments</li>
                <li>Cancel anytime</li>
                <li>Prices are in USD and may include taxes depending on your jurisdiction</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8.1. Payments and Merchant of Record</h2>
              <p>
                All payments are handled by our payment processor, who acts as Merchant of Record for all transactions. Our payment processor manages billing, taxes, invoicing, and refunds.
              </p>
              <p>
                MiCartaPro does not store or process payment information.
              </p>
              <p>
                Payment inquiries and refund requests must be addressed through our support channel at support@micartapro.com.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8.2. Cancellation and Refunds</h2>
              <p>
                You can cancel your subscription at any time when you find it convenient. There are no cancellation fees or penalties.
              </p>
              <p>
                When you cancel, your subscription will remain active until the end of your current billing period. Cancellation prevents future charges but does not generate refunds for past payments.
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>We do not offer a money-back guarantee</li>
                <li>Refund requests are handled at our discretion on a case-by-case basis</li>
                <li>No refunds for unused subscription time</li>
                <li>All refunds, if granted, are processed by our payment processor</li>
              </ul>
              <p className="mt-4">
                See also: <Link href="/refund" className="text-blue-600 hover:text-blue-700 underline">Refund Policy</Link>
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Service Availability</h2>
              <p>
                We aim for continuous availability, but service may experience interruptions due to maintenance or technical issues. We provide no uptime guarantee.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Limitation of Liability</h2>
              <p>
                To the fullest extent allowed by law, MiCartaPro is not liable for indirect, incidental, or consequential damages (including lost profits, data loss, or business interruptions).
              </p>
              <p>
                Maximum liability for any claim is limited to amounts paid in the last 12 months.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Indemnification</h2>
              <p>
                You agree to defend and indemnify MiCartaPro against claims arising from:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Your use of the Services</li>
                <li>Violation of these Terms</li>
                <li>Violation of laws or third-party rights</li>
                <li>Uploading User Content</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Termination</h2>
              <p>
                We may suspend or terminate accounts without notice if Terms are violated, fraudulent activity is detected, or fees remain unpaid.
              </p>
              <p className="mt-4">
                <strong>Cancellation by user:</strong>
              </p>
              <p>
                You may cancel anytime by contacting support at least 48 hours before renewal. Service remains active until the end of the billing period.
              </p>
              <p>
                Unused subscription time is non-refundable.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">13. Add-On Services (Optional)</h2>
              <p>
                Additional services may be offered separately, such as:
              </p>
              <ul className="list-disc pl-6 space-y-2 mb-4">
                <li>Custom menu design</li>
                <li>Manual menu setup</li>
                <li>Consultation</li>
                <li>Migration support</li>
              </ul>
              <p>These are:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Optional and not required</li>
                <li>Billed separately (bank transfer, PayPal, invoice)</li>
                <li>Not processed through our standard payment processor</li>
                <li>From $150 USD depending on complexity</li>
                <li>Separate terms may apply</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">14. Third-Party Services</h2>
              <p>
                The platform may integrate with third-party tools such as WhatsApp. Their terms and privacy policies apply separately.
              </p>
              <p>
                We are not responsible for third-party service behavior or data handling.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">15. Modifications to Terms</h2>
              <p>
                We may update these Terms periodically. The "Last Updated" date will change accordingly. Continued use of the Service implies acceptance of updated Terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">16. Governing Law and Jurisdiction</h2>
              <p>
                Except where consumer law requires otherwise, these Terms are governed by the laws of England and Wales.
              </p>
              <p>
                Any disputes shall be resolved in courts located in London, United Kingdom, unless arbitration or another method is mutually agreed.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">17. General Provisions</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">17.1. Entire Agreement</h3>
              <p>These Terms form the complete agreement between you and MiCartaPro.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">17.2. Severability</h3>
              <p>If any provision is deemed unenforceable, the remaining terms remain valid.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">17.3. Waiver</h3>
              <p>Failure to enforce rights under these Terms does not waive future rights.</p>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">17.4. Assignment</h3>
              <p>You may not transfer rights under these Terms without written consent. MiCartaPro may assign these Terms freely.</p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">18. Privacy & Cookies</h2>
              <p>
                Use of the Service is also governed by our Privacy Policy and Cookie Policy.
              </p>
              <p>
                By using MiCartaPro, you consent to data handling described in those documents.
              </p>
              <p className="mt-4">
                See: <Link href="/privacy" className="text-blue-600 hover:text-blue-700 underline">Privacy Policy</Link>
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">19. Contact</h2>
              <p>
                For questions regarding these Terms:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56 9 5785 7558</li>
                <li><strong>Email:</strong> legal@micartapro.com</li>
              </ul>
            </section>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-200 flex flex-col sm:flex-row justify-between items-center gap-4">
            <Link 
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              ← Back to Home
            </Link>
            <div className="flex flex-col sm:flex-row gap-4">
              <Link 
                href="/privacy"
                className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
              >
                View Privacy Policy →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}
