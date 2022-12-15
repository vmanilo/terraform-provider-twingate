// *** WARNING: this file was generated by the Pulumi Terraform Bridge (tfgen) Tool. ***
// *** Do not edit by hand unless you're certain you know what you are doing! ***

using System;
using System.Collections.Generic;
using System.Collections.Immutable;
using System.Threading.Tasks;
using Pulumi.Serialization;
using Pulumi;

namespace TwingateLabs.Twingate.Inputs
{

    public sealed class GetTwingateUsersUserArgs : global::Pulumi.InvokeArgs
    {
        [Input("email", required: true)]
        public string Email { get; set; } = null!;

        [Input("firstName", required: true)]
        public string FirstName { get; set; } = null!;

        [Input("id", required: true)]
        public string Id { get; set; } = null!;

        [Input("isAdmin", required: true)]
        public bool IsAdmin { get; set; }

        [Input("lastName", required: true)]
        public string LastName { get; set; } = null!;

        [Input("role", required: true)]
        public string Role { get; set; } = null!;

        public GetTwingateUsersUserArgs()
        {
        }
        public static new GetTwingateUsersUserArgs Empty => new GetTwingateUsersUserArgs();
    }
}